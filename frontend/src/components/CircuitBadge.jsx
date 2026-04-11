export function CircuitBadge({ state }) {
  const styles = {
    CLOSED:    { bg: "#16a34a", label: "CLOSED",    desc: "Traffic → Primary" },
    OPEN:      { bg: "#dc2626", label: "OPEN",      desc: "Traffic → Fallback" },
    HALF_OPEN: { bg: "#d97706", label: "HALF-OPEN", desc: "Probing Primary..." },
  };

  const s = styles[state] || styles["CLOSED"];

  return (
    <div style={{
      display: "inline-flex", flexDirection: "column", alignItems: "center",
      gap: 4, padding: "12px 24px",
      background: s.bg, borderRadius: 12, color: "#fff",
      minWidth: 160, textAlign: "center"
    }}>
      <span style={{ fontSize: 22, fontWeight: 700, letterSpacing: 1 }}>{s.label}</span>
      <span style={{ fontSize: 12, opacity: 0.85 }}>{s.desc}</span>
    </div>
  );
}