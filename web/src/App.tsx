import { useEffect, useState } from "react";
import "./App.css";

import { CPUSingleLineChart } from "@/components/CPUSingleLineChart";
import {
  Card,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export interface MemStat {
  totalMb: number;
  usedMb: number;
  availableMb: number;
  usagePercent: number;
  swapTotalMb: number;
  swapUsedMb: number;
  swapUsagePercent: number;
}

export interface AllData {
  cpu: Record<string, number>;
  mem: MemStat;
}

const COLOR_PALETTE = [
  "#3b82f6",
  "#10b981",
  "#f59e0b",
  "#ef4444",
  "#8b5cf6",
  "#ec4899",
  "#06b6d4",
  "#84cc16",
];

function App() {
  const [allStat, setAllStat] = useState<AllData | null>(null);
  const [history, setHistory] = useState<AllData[]>([]);

  useEffect(() => {
    const eventSource = new EventSource("http://127.0.0.1:8080/cpu-stat");

    eventSource.onmessage = (event) => {
      const data: AllData = JSON.parse(event.data);

      setAllStat(data);

      setHistory((prev) => [data, ...prev.slice(0, 29)]);
    };

    eventSource.onerror = () => {
      console.log("SSE disconnected, attempting to reconnect");
    };

    return () => eventSource.close();
  }, []);

  const cpuHistory = history.map((data) => data.cpu);
  const availableMetrics = allStat ? Object.keys(allStat.cpu) : [];

  if (!allStat) {
    return (
      <div className="text-amber-300">
        <h2>🚀 Connecting to Super Stat Go Stream...</h2>
      </div>
    );
  }

  return (
    <div className="py-6 ">
      <div className="border-b border-slate-300">
        <h1 className="text-slate-900 text-2xl font-bold tracking-wide mb-4 px-6">
          Dashboard
        </h1>
      </div>

      <Card className="ring-0 p-6">
        <CardHeader className="px-2">
          <CardTitle>CPU Usage</CardTitle>
        </CardHeader>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4 px-2">
          {availableMetrics.reverse().map((key, index) => (
            <CPUSingleLineChart
              key={key}
              history={cpuHistory}
              metricKey={key}
              color={COLOR_PALETTE[index % COLOR_PALETTE.length]}
            />
          ))}
        </div>
      </Card>
      <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl">
        <span className="text-xs text-slate-400 font-mono">RAM USAGE</span>
        <div className="text-3xl font-bold font-mono text-emerald-400 mt-1">
          {allStat.mem.usagePercent}%
        </div>
        <p className="text-xs text-slate-400 mt-1 font-mono">
          {allStat.mem.usedMb} MB / {allStat.mem.totalMb} MB
        </p>
      </div>
      <div className="bg-slate-900 border border-slate-800 p-4 rounded-xl">
        <span className="text-xs text-slate-400 font-mono">Swap USAGE</span>
        <div className="text-3xl font-bold font-mono text-emerald-400 mt-1">
          {allStat.mem.swapUsagePercent}%
        </div>
        <p className="text-xs text-slate-400 mt-1 font-mono">
          {allStat.mem.swapUsedMb} MB / {allStat.mem.swapTotalMb} MB
        </p>
      </div>
    </div>
  );
}

export default App;
