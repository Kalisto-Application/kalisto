import React, { useContext } from "react";
import { EditorSwitcher } from "../ui/EditorSwitcher";
import { CodeEditor } from "./CodeEditor";

import { Context } from "../state";
import { UpdateWorkspace } from "../../wailsjs/go/api/Api";
import { models } from "../../wailsjs/go/models";

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);

  const saveScript = (script: string) => {
    if (!ctx.state.activeWorkspaceId) {
      console.log("no active workspace");
      return;
    }

    const ws = ctx.state.workspaceList?.find(
      (it) => it.id === ctx.state.activeWorkspaceId,
    );
    if (!ws) {
      console.log("workspace not found");
      return;
    }
    const updatedWs = new models.Workspace({ ...ws, script: script });
    console.log(updatedWs);

    UpdateWorkspace(updatedWs)
      .then((_) => {
        console.log("workspace script saved");
      })
      .catch((err) => {
        console.log("failed to save script: ", err);
        if (err?.Code === "SYNTAX_ERROR") {
          ctx.dispatch({ type: "scriptError", value: err.Value });
          return;
        }
        console.log("failed to save global vars: ", err);
      });
  };

  return (
    <div className="bg-textBlockFill w-1/2">
      <CodeEditor
        key={0}
        text={ctx.state.scriptText}
        type="changeScriptText"
        action={saveScript}
      />
    </div>
  );
};
