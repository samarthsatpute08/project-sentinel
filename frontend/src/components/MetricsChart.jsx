import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Legend } from "recharts";

export function MetricsChart({ history }) {
  return (
    <div style={{ width: "100%", height: 200 }}>
      <ResponsiveContainer>
        <LineChart data={history} margin={{ top: 5, right: 20, left: 0, bottom: 5 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#374151"/>
          <XAxis dataKey="time" tick={{ fontSize: 10, fill: "#9ca3af" }} interval="preserveStartEnd"/>
          <YAxis tick={{ fontSize: 10, fill: "#9ca3af" }}/>
          <Tooltip
            contentStyle={{ background: "#1f2937", border: "1px solid #374151", borderRadius: 8 }}
            labelStyle={{ color: "#f9fafb" }}
          />
          <Legend/>
          <Line
            type="monotone" dataKey="rps" stroke="#3b82f6"
            dot={false} strokeWidth={2} name="RPS" isAnimationActive={false}
          />
          <Line
            type="monotone" dataKey="latency" stroke="#f97316"
            dot={false} strokeWidth={2} name="Latency (ms)" isAnimationActive={false}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}