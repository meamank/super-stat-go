import { AreaChart, Area, XAxis, YAxis, Tooltip, CartesianGrid } from "recharts"
import { ChartContainer, type ChartConfig } from "@/components/ui/chart"

    export function CPUSingleLineChart({
      history,
      metricKey,
      color = "#3b82f6"
    }: {
      history: Record<string, number>[];
      metricKey: string;
      color?: string;
    }) {
      const latestValue = history[0]?.[metricKey] ?? 0;
      const chartConfig = {
        [metricKey]: { label: metricKey.toUpperCase(), color: color },
      } satisfies ChartConfig

      const chartData = [...history].reverse().map((item, index) => ({
        time: `${history.length - index}s`,
        value: item[metricKey] ?? 0,
      }));

      return (
        <div className="bg-white border border-slate-300 p-4 rounded-xl">
          <h3 className="font-mono text-sm font-bold text-slate-900 mb-2">{metricKey.toUpperCase()} - {latestValue}%</h3>
          <ChartContainer config={chartConfig} className="h-32 w-full">
            <AreaChart data={chartData} margin={{ left: -10, right: 10, top: 5, bottom: 5 }}>
              <CartesianGrid vertical={false} />
              <XAxis dataKey="time" hide />
              <YAxis domain={[0, 100]} tickCount={4} tickLine={false} axisLine={false} hide/>
              <Tooltip />
              <Area type="monotone" dataKey="value" stroke={color} fill={color} fillOpacity={0.2} strokeWidth={2} />
            </AreaChart>
          </ChartContainer>
        </div>
      )
    }