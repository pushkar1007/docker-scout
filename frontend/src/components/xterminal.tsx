"use client";

import { Terminal } from "@xterm/xterm";
import { useEffect, useRef } from "react";
import "@xterm/xterm/css/xterm.css";
import { connectTerminalWS } from "./helpers/webSocket";

const WEBSOCKET =
  process.env.NEXT_PUBLIC_WEBSOCKET || "ws://localhost:8080/ws";



export default function XTerminal() {
  const terminalRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new Terminal({ cursorBlink: true, fontSize: 14 });

    const charWidth = 8;
    const charHeight = 16;
    const width = window.innerWidth;
    const height = window.innerHeight;
    const cols = Math.floor(width / charWidth);
    const rows = Math.floor(height / charHeight);

    term.resize(cols, rows);
    term.open(terminalRef.current);


    const ws = connectTerminalWS(
      WEBSOCKET,
      (data) => term.write('\x1B[1;3;31m :$> \x1B[0m' + data), // incoming PTY bytes
      () => term.writeln("\r\n\x1B[1;3;31m :$> \x1B[0mConnected to the terminal"),
      () => term.writeln("\r\n\x1B[1;3;31m :$> \x1B[0mDisconnected."),
      (err) => console.error("WS error:", err)
    );

    // send keystrokes to backend
    term.onData((e) => ws.send(e));

    // cleanup
    return () => {
      ws.close();
      term.dispose();
    };
  }, []);

  return <div ref={terminalRef} className="w-full h-screen bg-fuchsia-200" />;
}

