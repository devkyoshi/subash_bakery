import { Card } from "@/components/ui/card";
import { Calendar } from "@/components/ui/calendar";
import { useMemo, useState } from "react";

export function CalendarCard() {
  const [date, setDate] = useState<Date | undefined>(new Date("2021-05-10"));
  const defaultMonth = useMemo(() => new Date("2021-05-01"), []);

  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-center gap-2 text-sm font-semibold">
        Calendar
      </div>
      <div className="mt-4 flex justify-center">
        <Calendar
          mode="single"
          selected={date}
          onSelect={setDate}
          defaultMonth={defaultMonth}
          className="rounded-2xl border border-border bg-background p-4"
        />
      </div>
    </Card>
  );
}
