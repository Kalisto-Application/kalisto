import React, { useContext } from "react";
import { EditorSwitcher } from "../ui/EditorSwitcher";
import { CodeEditor } from "./CodeEditor";

import { Context } from "../state";

export const ScriptEditor: React.FC = () => {
    const ctx = useContext(Context)
    return (
        <div className="bg-textBlockFill w-1/2">
            <CodeEditor key={0} text={ctx.state.scriptText} type='changeScriptText' />
        </div>
    );
}
