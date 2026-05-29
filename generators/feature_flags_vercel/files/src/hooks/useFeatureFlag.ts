import { useState, useEffect } from "react";
import { getFlag, type FeatureFlag } from "../lib/flags";

export function useFeatureFlag(flag: FeatureFlag): boolean {
  const [enabled, setEnabled] = useState(false);
  useEffect(() => {
    getFlag(flag).then(setEnabled);
  }, [flag]);
  return enabled;
}
