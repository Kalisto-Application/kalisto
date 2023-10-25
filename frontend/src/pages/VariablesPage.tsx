import React, { useContext } from 'react';
import { SaveGlovalVars } from '../../wailsjs/go/api/Api';
import VarsError from '../components/VarsError';
import { Context } from '../state';
import { CodeEditor } from '../ui/Editor';

type VariablesPageProps = {};

export const VariablesPage: React.FC<VariablesPageProps> = () => {
  const ctx = useContext(Context);

  const saveGlobalVariables = (vars: string) => {
    SaveGlovalVars(vars)
      .then(() => {
        ctx.dispatch({ type: 'changeVariables', text: vars });
      })
      .catch((err) => {
        console.log('failed to save global vars: ', err);
        if (err?.Code === 'SYNTAX_ERROR') {
          ctx.dispatch({ type: 'varsError', value: err.Value });
          return;
        }
        console.log('failed to save global vars: ', err);
      });
  };

  return (
    <div className="flex flex-1">
      <div className="w-1/2 bg-textBlockFill">
        <CodeEditor
          fileId="globalVars"
          text={ctx.state.vars}
          onChange={saveGlobalVariables}
        />
        <div className="w-1/2"></div>
      </div>
      <VarsError />
    </div>
  );
};
