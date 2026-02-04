import { DashboardData } from "@/types/types";

export function connectTerminalWS(
  url: string,
  onMessage: (data: DashboardData) => void,
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
    onMessage(event.data);
  };

  ws.onclose = () => {
    if (onClose) onClose();
  };

  ws.onerror = (ev) => {
    if (onError) onError(ev);
  };

  return ws;
}
