import { useEffect, useRef, useState, useCallback } from "react";

const WS_URL = import.meta.env.VITE_WS_URL || "ws://localhost:8080/ws";
const RENDER_INTERVAL_MS = 300;
const MAX_HISTORY = 60; // 60 data points in chart

export function useMetrics() {
  const [metrics, setMetrics] = useState(null);
  const [history, setHistory] = useState([]);
  const [connected, setConnected] = useState(false);
  const bufferRef = useRef([]); // raw incoming messages, never triggers render
  const wsRef = useRef(null);

  const connect = useCallback(() => {
    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => {
      setConnected(false);
      // Auto-reconnect after 2s
      setTimeout(connect, 2000);
    };
    ws.onerror = () => ws.close();

    ws.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        bufferRef.current.push(data);
        // Keep buffer small to avoid memory growth
        if (bufferRef.current.length > 10) {
          bufferRef.current = bufferRef.current.slice(-5);
        }
      } catch {
        // ignore malformed messages
      }
    };
  }, []);

  // Connect on mount
  useEffect(() => {
    connect();
    return () => wsRef.current?.close();
  }, [connect]);

  // Drain buffer into state on interval — this is the key trick
  useEffect(() => {
    const interval = setInterval(() => {
      if (bufferRef.current.length === 0) return;

      // Take the latest message from the buffer
      const latest = bufferRef.current[bufferRef.current.length - 1];
      bufferRef.current = [];

      setMetrics(latest);
      setHistory((prev) => {
        const next = [...prev, { ...latest, time: new Date(latest.timestamp).toLocaleTimeString() }];
        return next.slice(-MAX_HISTORY);
      });
    }, RENDER_INTERVAL_MS);

    return () => clearInterval(interval);
  }, []);

  return { metrics, history, connected };
}