import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useThemeToggle } from "@/theme/useThemeToggle";

export function ThemeToggle() {
  const { theme, toggle } = useThemeToggle();
  const Icon = theme === "dark" ? Sun : Moon;

  return (
    <Button variant="ghost" size="icon" onClick={toggle} aria-label="Toggle theme">
      <Icon className="h-4 w-4" />
    </Button>
  );
}
