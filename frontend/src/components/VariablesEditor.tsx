import React, { useContext } from 'react';
import { SaveGlovalVars } from '../../wailsjs/go/api/Api';
import { Context } from '../state';
import { CodeEditor } from '../ui/Editor';

export const VariablesEditor: React.FC = () => {
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
    <CodeEditor
      fileId="globalVars"
      text={ctx.state.vars}
      onChange={saveGlobalVariables}
    />
  );
};
