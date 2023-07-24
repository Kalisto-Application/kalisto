import React, { useContext } from "react";
import { CodeEditor } from "../components/CodeEditor";
import { Context } from "../state";

type VariablesPageProps = {

}



export const VariablesPage: React.FC<VariablesPageProps> = ({}) => {
  const ctx = useContext(Context);

    return (
      <div className="p-4">
        <CodeEditor text={ctx.state.variables} type='changeVariables' action={(t: string) => {console.log(t)}} />
      </div>
    );
  };
