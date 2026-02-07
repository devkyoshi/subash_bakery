import { Card } from "@/components/ui/card";
import type { SalesPoint } from "@/types/dashboard";
import {
  Bar,
  ComposedChart,
  Line,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

function formatK(v: number) {
  return `${(v / 1000).toFixed(0)},000`;
}

export function SalesHistoryCard({ data }: { data: SalesPoint[] }) {
  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <div className="text-sm font-semibold">Sales history</div>
        </div>
        <div className="text-xs text-muted-foreground">Period: 1 year</div>
      </div>

      <div className="mt-4 h-[280px] w-full">
        <ResponsiveContainer width="100%" height="100%">
          <ComposedChart
            data={data}
            margin={{ left: 10, right: 10, top: 10, bottom: 0 }}
          >
            <XAxis dataKey="month" tick={{ fontSize: 11 }} interval={2} />
            <YAxis
              tick={{ fontSize: 11 }}
              tickFormatter={(v) => formatK(v as number)}
              width={50}
            />
            <Tooltip
              contentStyle={{
                borderRadius: 12,
                border: "1px solid hsl(var(--border))",
                background: "hsl(var(--elevated))",
              }}
              labelStyle={{ fontWeight: 600 }}
              formatter={(value: unknown, name: string) => {
                if (name === "cashedPct")
                  return [
                    `${Number(value).toFixed(2)}%`,
                    "Total cashed in (%)",
                  ];
                if (name === "cashed")
                  return [
                    Number(value).toLocaleString(),
                    "Total cashed in (Rs.)",
                  ];
                return [Number(value).toLocaleString(), "Total invoiced"];
              }}
            />
            <Bar
              dataKey="invoiced"
              fill="hsl(var(--info))"
              radius={[6, 6, 0, 0]}
              barSize={14}
            />
            <Bar
              dataKey="cashed"
              fill="hsl(var(--success))"
              radius={[6, 6, 0, 0]}
              barSize={14}
              opacity={0.9}
            />
            <Line
              type="monotone"
              dataKey="cashedPct"
              stroke="hsl(var(--success))"
              strokeWidth={2}
              dot={false}
            />
          </ComposedChart>
        </ResponsiveContainer>
      </div>
    </Card>
  );
}
