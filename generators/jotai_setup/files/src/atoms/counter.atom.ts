import { atom } from "jotai";

export const countAtom = atom(0);

export const incrementAtom = atom(null, (get, set) =>
  set(countAtom, get(countAtom) + 1),
);

export const decrementAtom = atom(null, (get, set) =>
  set(countAtom, get(countAtom) - 1),
);
