import React, { useContext } from "react";
import { Context } from "../state/state";

interface EditorSwitcherProps {
    switches: string[];
}

export const EditorSwitcher: React.FC<EditorSwitcherProps> = ({ switches }) => {
    const ctx = useContext(Context)
    const makeClassName = (i: number): string => {
        return ctx.state.activeEditor === i? "pt-[8px] px-[16px]": "pt-[8px] px-[16px] text-secondaryText"
    }

    return (
        <div className="flex flex-1 justify-start h-[40px] ">
            { switches.map((it, i) => <div className={makeClassName(i)} key={i} onClick={() => ctx.dispatch({type: 'switchEditor', i: i})}>{it}</div>) }
        </div>
    );
}
