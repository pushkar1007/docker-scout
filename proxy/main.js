const http = require("http");
const express = require("express");
const Docker = require("dockerode");
const httpProxy = require("http-proxy");

const docker = new Docker({ socketPath: "/var/run/docker.sock" });
const proxy = httpProxy.createProxy({ ws: true });

const managementAPI = express();
managementAPI.use(express.json());

const reverseProxyApp = express();
const db = new Map();

console.log("Listening for Docker container events...");

/* -------------------- DOCKER EVENTS -------------------- */

let buffer = "";

docker.getEvents((err, stream) => {
    if (err) {
        console.error("Docker events error:", err);
        process.exit(1);
    }

    stream.on("data", async (chunk) => {
        buffer += chunk.toString();

        let boundary;
        while ((boundary = buffer.indexOf("\n")) !== -1) {
            const line = buffer.slice(0, boundary).trim();
            buffer = buffer.slice(boundary + 1);

            if (!line) continue;

            try {
                const event = JSON.parse(line);

                if (event.Type !== "container") continue;

                const container = docker.getContainer(event.Actor.ID);

                if (event.Action === "start") {
                    const info = await container.inspect();

                    const name = info.Name.replace("/", "").toLowerCase();
                    const port = 3000;

                    db.set(name, { name, port });

                    console.log(
                        `Registered: ${name}.localhost → http://${name}:${port}`
                    );
                }

                if (event.Action === "die" || event.Action === "destroy") {
                    for (const [key] of db) {
                        if (key === event.Actor.Attributes.name) {
                            db.delete(key);
                            console.log(`Removed: ${key}`);
                        }
                    }
                }
            } catch (e) {
                console.error("Event parse error:", e.message);
            }
        }
    });
});

/* -------------------- REVERSE PROXY -------------------- */

reverseProxyApp.use((req, res) => {
    // URL: /container-name/anything
    const parts = req.url.split("/").filter(Boolean);

    if (parts.length === 0) {
        return res.status(200).send("Reverse proxy running");
    }

    const containerName = parts[0].toLowerCase();

    if (!db.has(containerName)) {
        return res.status(404).send("Container not found");
    }

    const { name, port } = db.get(containerName);
    const target = `http://${name}:${port}`;

    // Strip /container-name from path
    req.url = "/" + parts.slice(1).join("/");

    console.log(`HTTP /${containerName} → ${target}${req.url}`);

    proxy.web(req, res, { target });
});


const reverseProxy = http.createServer(reverseProxyApp);

reverseProxy.on("upgrade", (req, socket, head) => {
    const parts = req.url.split("/").filter(Boolean);

    if (parts.length === 0) {
        socket.destroy();
        return;
    }

    const containerName = parts[0].toLowerCase();

    if (!db.has(containerName)) {
        socket.destroy();
        return;
    }

    const { name, port } = db.get(containerName);
    const target = `http://${name}:${port}`;

    req.url = "/" + parts.slice(1).join("/");

    console.log(`WS /${containerName} → ${target}${req.url}`);

    proxy.ws(req, socket, head, { target });
});

/* -------------------- MANAGEMENT API -------------------- */

managementAPI.post("/containers", async (req, res) => {
    const { image, tag = "latest", name } = req.body;

    if (!image || !name) {
        return res.status(400).json({ error: "image and name required" });
    }

    const images = await docker.listImages();
    const exists = images.some(img =>
        (img.RepoTags || []).includes(`${image}:${tag}`)
    );

    if (!exists) {
        console.log(`Pulling ${image}:${tag}`);
        await new Promise((resolve, reject) => {
            docker.pull(`${image}:${tag}`, (err, stream) => {
                if (err) return reject(err);
                docker.modem.followProgress(stream, resolve);
            });
        });
    }

    const container = await docker.createContainer({
        Image: `${image}:${tag}`,
        name,
        HostConfig: {
            AutoRemove: true,
            NetworkMode: "proxy-net"
        }
    });

    await container.start();

    res.json({
        success: true,
        url: `http://${name}.localhost`
    });
});

/* -------------------- START SERVERS -------------------- */

managementAPI.listen(8080, () =>
    console.log("Management API running on :8080")
);

reverseProxy.listen(80, () =>
    console.log("Reverse Proxy running on :80")
);
