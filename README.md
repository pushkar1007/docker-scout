Here are all the API requests for the Docker Scout application:

## System Management
- `GET /system/health` - Health check
- `POST /system/nuke?confirm=true` - System cleanup (removes all containers, images, unused volumes)

## Containers
- `GET /containers` - List all containers
- `POST /containers/start?id=<container_id>` - Start container
- `POST /containers/stop?id=<container_id>` - Stop container
- `POST /containers/pause?id=<container_id>` - Pause container
- `POST /containers/unpause?id=<container_id>` - Unpause container
- `DELETE /containers/remove?id=<container_id>` - Remove container

## Images
- `GET /images` - List all images
- `DELETE /images/remove?name=<image_name>&force=true` - Remove image
- `POST /images/build` - Build image (with context_path and tag in JSON body)

## Volumes
- `GET /volumes` - List all volumes
- `GET /volumes/inspect?name=<volume_name>` - Inspect volume
- `POST /volumes/create` - Create volume (with name and driver in JSON body)
- `DELETE /volumes/remove?name=<volume_name>` - Remove volume
- `POST /volumes/prune` - Prune unused volumes

## Statistics & Monitoring
- `GET /stats` - Get container statistics
- `GET /events` - Server-Sent Events stream (real-time Docker events)
- `POST /events/publish` - Publish custom event (with event and data in JSON body)

**Server:** Runs on `http://localhost:8089`