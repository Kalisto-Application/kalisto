import React, { useContext } from "react";
import { Context } from "../state/state";

interface EditorSwitcherProps {
    switches: string[];
}

export const EditorSwitcher: React.FC<EditorSwitcherProps> = ({ switches }) => {
    const ctx = useContext(Context)

    return (
        <div>
            { switches.map((it, i) => <span key={i} onClick={() => ctx.dispatch({type: 'switchEditor', i: i})}>{it}</span>) }
        </div>
    );
}
