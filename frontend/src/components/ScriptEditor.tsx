import React, { useContext, useEffect, useState } from 'react';
import { CodeEditor } from './CodeEditor2';

import { Context } from '../state';
import { UpdateWorkspace } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);

  let idFile = ctx.state.scriptIdFile;
  const ws = ctx.state.activeWorkspace;

  const activeFile = ws?.scriptFiles?.find((it) => it.id == idFile);

  const saveScript = (script: string) => {
    debugger;
    // console.log('value script', script);

    if (!ws) {
      console.log('no active workspace');
      return;
    }

    if (!ws) {
      console.log('workspace not found');
      return;
    }

    const updatedWs = new models.Workspace({
      ...ws,
      scriptFiles: ws.scriptFiles?.map((file) => {
        if (idFile === file.id) {
          file.content = script;
          return file;
        }
        return file;
      }),
    });

    UpdateWorkspace(updatedWs)
      .then((_) => {
        ctx.dispatch({
          type: 'updateWorkspace',
          workspace: updatedWs,
        });
        console.log('workspace script saved');
      })
      .catch((err) => {
        console.log('failed to save script: ', err);
        if (err?.Code === 'SYNTAX_ERROR') {
          ctx.dispatch({ type: 'scriptError', value: err.Value });
          return;
        }
        console.log('failed to save global vars: ', err);
      });
  };
  console.log('scriptFiles update', ctx.state.activeWorkspace?.scriptFiles);

  return (
    <div className="w-1/2 bg-textBlockFill">
      <CodeEditor
        text={activeFile?.content || ''}
        action={saveScript}
        idFile={idFile}
      />
    </div>
  );
};
