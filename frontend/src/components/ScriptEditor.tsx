import React, { useContext, useEffect, useState } from 'react';
import { CodeEditor } from '../ui/Editor';

import { Context } from '../state';
import { UpdateWorkspace } from '../../wailsjs/go/api/Api';
import { models } from '../../wailsjs/go/models';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);

  let idFile = ctx.state.scriptIdFile;
  const ws = ctx.state.activeWorkspace;

  const activeFile = ws?.scriptFiles?.find((it) => it.id == idFile);

  const saveScript = (content: string) => {
    ctx.dispatch({
      type: 'updateScriptFile',
      content: content,
    });
    console.log('workspace script saved');
  };

  return (
    <div className="w-1/2 bg-textBlockFill">
      <CodeEditor
        text={activeFile?.content || ''}
        fileId={ctx.state.scriptIdFile}
        onChange={saveScript}
      />
    </div>
  );
};
