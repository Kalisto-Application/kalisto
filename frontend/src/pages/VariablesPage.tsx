import React, { useContext } from "react";
import { CodeEditor } from "../components/CodeEditor";
import { Context } from "../state";
import { SaveGlovalVars } from "../../wailsjs/go/api/Api";

type VariablesPageProps = {
}

export const VariablesPage: React.FC<VariablesPageProps> = () => {
  const ctx = useContext(Context);

    const saveGlobalVariables = (vars: string) => {
      SaveGlovalVars(vars).catch(err => {
        console.log('failed to save global vars: ', err)
      })
    }

    return (
      <div className="">
        <CodeEditor text={ctx.state.vars} type='changeVariables' action={saveGlobalVariables} />
      </div>
    );
  };
