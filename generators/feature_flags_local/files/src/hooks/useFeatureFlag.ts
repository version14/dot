import { useState, useEffect } from "react";
import { getFlag } from "../lib/flags";

export function useFeatureFlag(flag: string): boolean {
  const [enabled, setEnabled] = useState(false);
  useEffect(() => {
    getFlag(flag).then(setEnabled);
  }, [flag]);
  return enabled;
}
