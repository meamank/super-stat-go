import { Tooltip, PieChart, Pie, Cell } from "recharts";
import { ChartContainer, type ChartConfig } from "@/components/ui/chart";

interface MemPieChartProps {
  title: string;
  totalMB: number;
  usedMB: number;
  availableMB: number;
  usagePercent: number;
  usedColor?: string;
  availableColor?: string;
}

export function MemDonutChart({
  title,
  totalMB,
  usedMB,
  availableMB,
  usagePercent,
  usedColor = "#ef4444",
  availableColor = "#10b981",
}: MemPieChartProps) {
  const chartConfig = {
    used: { label: "Used", color: usedColor },
    available: { label: "Available", color: availableColor },
  } satisfies ChartConfig;

  const chartData = [
    { name: "Used", value: usedMB, fill: usedColor },
    { name: "Available", value: availableMB, fill: availableColor },
  ];

  return (
    <div
      className="flex flex-col items-center justify-center p-4 bg-white border border-slate-200 rounded-xl
  shadow-sm"
    >
      <h3 className="self-start font-mono text-sm font-semibold text-slate-600 mb-1">
        {title} - {totalMB} MB
      </h3>

      <div className="relative w-full h-32 flex items-center justify-center">
        <ChartContainer config={chartConfig} className="h-32 w-full">
          <PieChart>
            <Tooltip formatter={(value: number) => [`${value} MB`, "Memory"]} />
            <Pie
              data={chartData}
              dataKey="value"
              nameKey="name"
              innerRadius={40}
              outerRadius={60}
              paddingAngle={4}
              cornerRadius={6}
            >
              {chartData.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.fill} />
              ))}
            </Pie>
          </PieChart>
        </ChartContainer>

        {/* Center Text displaying percentage */}
        <div className="absolute flex flex-col items-center justify-center text-center font-mono">
          <span className="text-lg font-extrabold text-slate-900">
            {usagePercent}%
          </span>
          <span className="text-xs text-slate-500 uppercase tracking-wider">
            Used
          </span>
        </div>
      </div>

      {/* Legend below donut */}
      <div className="flex items-center gap-6 mt-1 font-mono text-xs">
        <div className="flex items-center gap-2">
          <span
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: usedColor }}
          ></span>
          <span>Used: {usedMB} MB</span>
        </div>
        <div className="flex items-center gap-2">
          <span
            className="w-3 h-3 rounded-full"
            style={{ backgroundColor: availableColor }}
          ></span>
          <span>Free: {availableMB} MB</span>
        </div>
      </div>
    </div>
  );
}
