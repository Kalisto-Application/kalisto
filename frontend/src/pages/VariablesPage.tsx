import React, { useContext } from "react";
import { CodeEditor } from "../components/CodeEditor";
import { Context } from "../state";
import { SaveGlovalVars } from "../../wailsjs/go/api/Api";
import VarsError from "../components/VarsError";

type VariablesPageProps = {};

export const VariablesPage: React.FC<VariablesPageProps> = () => {
  const ctx = useContext(Context);

  const saveGlobalVariables = (vars: string) => {
    SaveGlovalVars(vars).catch((err) => {
      console.log("failed to save global vars: ", err);
      if (err?.Code === "SYNTAX_ERROR") {
        console.log("dispatcyhed  ");
        ctx.dispatch({ type: "varsError", value: err.Value });
        return;
      }
      console.log("failed to save global vars: ", err);
    });
  };

  return (
    <div className="flex flex-1">
      <div className="bg-textBlockFill w-1/2">
        <CodeEditor
          text={ctx.state.vars}
          type="changeVariables"
          action={saveGlobalVariables}
        />
        <div className="w-1/2"></div>
      </div>
      <VarsError />
    </div>
  );
};
