import { Link } from "react-router-dom";
import { LogoMark } from "@/components/common/LogoMark";
import { Button } from "@/components/ui/button";

export function WelcomePage() {
  return (
    <div className="grid min-h-screen grid-cols-2">
      <div className="flex flex-col items-center justify-center gap-6 bg-background">
        <LogoMark />
        <div className="text-center">
          <div className="text-5xl font-semibold tracking-tight">Welcome!</div>
          <div className="mt-4 text-xl text-muted-foreground">Sign in to continue to the services</div>
        </div>
      </div>

      <div className="flex items-center justify-center border-l border-border bg-background">
        <div className="w-[420px]">
          <div className="text-center text-xl font-medium">Get started</div>
          <div className="mt-10 grid gap-4">
            <Button asChild className="h-12 rounded-none" variant="outline">
              <Link to="/auth/login">Sign in</Link>
            </Button>
            <Button asChild className="h-12 rounded-none" variant="outline">
              <Link to="/auth/register">Register new user</Link>
            </Button>
          </div>
          <div className="mt-28 text-center text-xs text-muted-foreground">
            By signing in you agree to our<br />
            <span className="text-foreground">Terms of Service</span> and <span className="text-foreground">Privacy Policy</span>
          </div>
        </div>
      </div>
    </div>
  );
}
