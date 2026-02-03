import { useTheme } from "next-themes";

export function useThemeToggle() {
  const { theme, setTheme, resolvedTheme } = useTheme();
  const current = (resolvedTheme ?? theme) as "light" | "dark" | undefined;

  return {
    theme: current,
    toggle: () => setTheme(current === "dark" ? "light" : "dark"),
    setTheme,
  };
}
