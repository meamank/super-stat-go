import { useEffect, useState } from "react";
import "./App.css";

import { CPUSingleLineChart } from "@/components/CPUSingleLineChart";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";

import { MemDonutChart } from "@/components/MemDonutChart";

export interface MemStat {
  [key: string]: number;
  totalMB: number;
  usedMB: number;
  availableMB: number;
  usagePercent: number;
  swapTotalMB: number;
  swapUsedMB: number;
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
    const eventSource = new EventSource("/cpu-stat");

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
  const memHistory = history.map((data) => data.mem);

  const availableMetricsCPU = allStat ? Object.keys(allStat.cpu) : [];

  console.log("cpu:", cpuHistory);
  console.log("mem", memHistory);

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
          {availableMetricsCPU.reverse().map((key, index) => (
            <CPUSingleLineChart
              key={key}
              history={cpuHistory}
              metricKey={key}
              color={COLOR_PALETTE[index % COLOR_PALETTE.length]}
            />
          ))}
        </div>
      </Card>

      {/* 2. Memory & Swap Telemetry Info */}
      <Card className="ring-0 p-6">
        <CardHeader className="px-2">
          <CardTitle>Memory & Swap Usage</CardTitle>
        </CardHeader>
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-6 px-2">
          <MemDonutChart
            title="Total RAM"
            totalMB={allStat.mem.totalMB}
            usedMB={allStat.mem.usedMB}
            availableMB={allStat.mem.availableMB}
            usagePercent={allStat.mem.usagePercent}
            usedColor="#ef4444"
            availableColor="#10b981"
          />

          {/* 2. Swap Memory Pie Chart */}
          <MemDonutChart
            title="Swap Memory"
            totalMB={allStat.mem.swapTotalMB}
            usedMB={allStat.mem.swapUsedMB}
            availableMB={allStat.mem.swapTotalMB - allStat.mem.swapUsedMB}
            usagePercent={allStat.mem.swapUsagePercent}
            usedColor="#f59e0b"
            availableColor="#3b82f6"
          />
        </div>
      </Card>
    </div>
  );
}

export default App;
