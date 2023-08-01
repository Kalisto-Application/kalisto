import React, { useContext } from "react";
import { Context } from "../state/state";

interface props {
  items: item[];
  active: number;
}

type item = {
  title: string;
  onClick: (i: number) => void;
}

export const EditorSwitcher: React.FC<props> = ({ items, active }) => {
    const makeClassName = (i: number): string => {
        return active === i? "pt-[8px] px-[16px]": "pt-[8px] px-[16px] text-secondaryText"
    }

    return (
        <div className="flex flex-1 justify-start h-[40px] ">
            { items.map((it, i) => <div className={makeClassName(i)} key={i} onClick={() => it.onClick(i)}>{it.title}</div>) }
        </div>
    );
}
