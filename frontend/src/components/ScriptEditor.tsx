import React, { useContext } from "react";
import { EditorSwitcher } from "../ui/EditorSwitcher";
import { CodeEditor } from "./CodeEditor";

import { Context } from "../state";
import { UpdateWorkspace, SaveScript } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

export const ScriptEditor: React.FC = () => {
    const ctx = useContext(Context)

    const saveScript = (script: string) => {
        if (!ctx.state.activeWorkspaceId) return;

        SaveScript(ctx.state.activeWorkspaceId, script).catch(err => {
          console.log('failed to save script: ', err)
          if (err?.Code === "SYNTAX_ERROR") {
            ctx.dispatch({type: 'scriptError', value: err.Value})
            return
          }
          console.log('failed to save global vars: ', err)
        })
      }

    return (
        <div className="bg-textBlockFill w-1/2">
            <CodeEditor key={0} text={ctx.state.scriptText} type='changeScriptText' action={saveScript} />
        </div>
    );
}
