export function TrafficFlow({ activeRoute, circuitState }) {
  const primaryActive = activeRoute === "primary";
  const fallbackActive = activeRoute === "fallback";

  const colors = {
    active:   "#2563eb",
    inactive: "#9ca3af",
    danger:   "#dc2626",
  };

  return (
    <div style={{ padding: "16px 0" }}>
      <svg width="100%" viewBox="0 0 500 160" style={{ overflow: "visible" }}>
        <defs>
          <marker id="a" viewBox="0 0 10 10" refX="8" refY="5"
            markerWidth="5" markerHeight="5" orient="auto-start-reverse">
            <path d="M2 1L8 5L2 9" fill="none" stroke="context-stroke"
              strokeWidth="1.5" strokeLinecap="round"/>
          </marker>
        </defs>

        {/* Client box */}
        <rect x="10" y="60" width="80" height="40" rx="8"
          fill="#1e40af" stroke="#3b82f6" strokeWidth="1"/>
        <text x="50" y="84" textAnchor="middle" fill="#fff" fontSize="13" fontWeight="600">Client</text>

        {/* Router box */}
        <rect x="200" y="55" width="100" height="50" rx="8"
          fill="#4c1d95" stroke="#8b5cf6" strokeWidth="1"/>
        <text x="250" y="78" textAnchor="middle" fill="#fff" fontSize="13" fontWeight="600">Router</text>
        <text x="250" y="95" textAnchor="middle" fill="#c4b5fd" fontSize="10">Circuit Breaker</text>

        {/* Primary API box */}
        <rect x="400" y="20" width="90" height="40" rx="8"
          fill={primaryActive ? "#166534" : "#374151"}
          stroke={primaryActive ? "#22c55e" : "#6b7280"} strokeWidth="1"/>
        <text x="445" y="44" textAnchor="middle" fill="#fff" fontSize="12" fontWeight="600">Primary</text>

        {/* Fallback API box */}
        <rect x="400" y="100" width="90" height="40" rx="8"
          fill={fallbackActive ? "#7c2d12" : "#374151"}
          stroke={fallbackActive ? "#f97316" : "#6b7280"} strokeWidth="1"/>
        <text x="445" y="124" textAnchor="middle" fill="#fff" fontSize="12" fontWeight="600">Fallback</text>

        {/* Client → Router */}
        <line x1="90" y1="80" x2="198" y2="80"
          stroke={colors.active} strokeWidth="2" markerEnd="url(#a)"/>

        {/* Router → Primary */}
        <line x1="300" y1="72" x2="398" y2="48"
          stroke={primaryActive ? colors.active : colors.inactive}
          strokeWidth={primaryActive ? 3 : 1.5}
          strokeDasharray={primaryActive ? "none" : "5,4"}
          markerEnd="url(#a)"/>

        {/* Router → Fallback */}
        <line x1="300" y1="88" x2="398" y2="112"
          stroke={fallbackActive ? colors.danger : colors.inactive}
          strokeWidth={fallbackActive ? 3 : 1.5}
          strokeDasharray={fallbackActive ? "none" : "5,4"}
          markerEnd="url(#a)"/>

        {/* Pulsing dot on active route */}
        {primaryActive && (
          <circle r="5" fill="#22c55e" opacity="0.9">
            <animateMotion dur="0.8s" repeatCount="indefinite"
              path="M90,80 L398,48"/>
          </circle>
        )}
        {fallbackActive && (
          <circle r="5" fill="#f97316" opacity="0.9">
            <animateMotion dur="0.8s" repeatCount="indefinite"
              path="M90,80 L398,112"/>
          </circle>
        )}
      </svg>
    </div>
  );
}