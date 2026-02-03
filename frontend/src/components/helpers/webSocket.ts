export function connectTerminalWS(
  url: string,
  onMessage: (data: Uint8Array) => void,
  onOpen?: () => void,
  onClose?: () => void,
  onError?: (ev: Event) => void,
): WebSocket {
  const ws = new WebSocket(url);
  ws.binaryType = "arraybuffer";

  ws.onopen = () => {
    if (onOpen) onOpen();
  };

  ws.onmessage = (event) => {
    if (event.data instanceof ArrayBuffer) {
      onMessage(new Uint8Array(event.data));
    } else if (typeof event.data === "string") {
      onMessage(new TextEncoder().encode(event.data));
    }
  };

  ws.onclose = () => {
    if (onClose) onClose();
  };

  ws.onerror = (ev) => {
    if (onError) onError(ev);
  };

  return ws;
}
