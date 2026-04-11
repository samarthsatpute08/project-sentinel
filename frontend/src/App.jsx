import { useMetrics } from "./hooks/useMetrics";
import { CircuitBadge } from "./components/CircuitBadge";
import { TrafficFlow } from "./components/TrafficFlow";
import { MetricsChart } from "./components/MetricsChart";

function StatCard({ label, value, color }) {
  return (
    <div style={{
      background: "#1f2937", borderRadius: 10, padding: "16px 20px",
      minWidth: 120, textAlign: "center", border: "1px solid #374151"
    }}>
      <div style={{ fontSize: 26, fontWeight: 700, color: color || "#f9fafb" }}>{value ?? "—"}</div>
      <div style={{ fontSize: 12, color: "#9ca3af", marginTop: 4 }}>{label}</div>
    </div>
  );
}

export default function App() {
  const { metrics, history, connected } = useMetrics();

  return (
    <div style={{
      minHeight: "100vh", background: "#111827", color: "#f9fafb",
      fontFamily: "system-ui, sans-serif", padding: 24
    }}>
      {/* Header */}
      <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 24 }}>
        <div style={{
          width: 10, height: 10, borderRadius: "50%",
          background: connected ? "#22c55e" : "#ef4444"
        }}/>
        <h1 style={{ margin: 0, fontSize: 22, fontWeight: 700 }}>
          Project Sentinel — War Room
        </h1>
        <span style={{ fontSize: 12, color: "#9ca3af" }}>
          {connected ? "Live" : "Connecting..."}
        </span>
      </div>

      {/* Circuit State */}
      <div style={{ marginBottom: 24 }}>
        <h2 style={{ fontSize: 14, color: "#9ca3af", marginBottom: 10 }}>Circuit Breaker State</h2>
        <CircuitBadge state={metrics?.circuitState || "CLOSED"}/>
      </div>

      {/* Stat Cards */}
      <div style={{ display: "flex", gap: 12, flexWrap: "wrap", marginBottom: 24 }}>
        <StatCard label="Requests/sec" value={metrics?.rps?.toFixed(1)} color="#3b82f6"/>
        <StatCard label="Total Requests" value={metrics?.totalRequests} color="#f9fafb"/>
        <StatCard label="Total Failures" value={metrics?.totalFailures} color="#ef4444"/>
        <StatCard label="Latency (ms)" value={metrics?.latency?.toFixed(0)} color="#f97316"/>
        <StatCard label="Active Route" value={metrics?.activeRoute} 
          color={metrics?.activeRoute === "primary" ? "#22c55e" : "#f97316"}/>
      </div>

      {/* Traffic Flow Diagram */}
      <div style={{ background: "#1f2937", borderRadius: 12, padding: 20, marginBottom: 24, border: "1px solid #374151" }}>
        <h2 style={{ fontSize: 14, color: "#9ca3af", margin: "0 0 8px" }}>Live Traffic Routing</h2>
        <TrafficFlow
          activeRoute={metrics?.activeRoute || "primary"}
          circuitState={metrics?.circuitState || "CLOSED"}
        />
      </div>

      {/* Metrics Chart */}
      <div style={{ background: "#1f2937", borderRadius: 12, padding: 20, border: "1px solid #374151" }}>
        <h2 style={{ fontSize: 14, color: "#9ca3af", margin: "0 0 12px" }}>RPS & Latency History</h2>
        <MetricsChart history={history}/>
      </div>
    </div>
  );
}