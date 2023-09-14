import React, { useContext } from "react";
import { Context } from "../state/state";

interface props {
  items: item[];
  active: number;
}

type item = {
  title: string;
  onClick: (i: number) => void;
};

export const EditorSwitcher: React.FC<props> = ({ items, active }) => {
  const makeClassName = (i: number): string => {
    return active === i
      ? "pt-[8px] px-[16px]"
      : "pt-[8px] px-[16px] text-secondaryText ";
  };

  return (
    <div className="flex h-[40px] flex-1 justify-start ">
      {items.map((it, i) => (
        <div
          className={`${makeClassName(i)} cursor-pointer`}
          key={i}
          onClick={() => it.onClick(i)}
        >
          {it.title}
        </div>
      ))}
    </div>
  );
};
