import React, { useContext, useEffect, useState } from 'react';
import { CodeEditor } from '../ui/Editor';

import { Context } from '../state';
import { UpdateScriptFileContent } from '../../wailsjs/go/api/Api';

export const ScriptEditor: React.FC = () => {
  const ctx = useContext(Context);

  let fileId = ctx.state.scriptIdFile;
  const ws = ctx.state.activeWorkspace;
  const activeFile = ws?.scriptFiles?.find((it) => it.id == fileId);

  const saveScript = (content: string) => {
    UpdateScriptFileContent(
      ctx.state.activeWorkspace?.id || '',
      fileId,
      content
    ).then((res) => {
      ctx.dispatch({
        type: 'updateScriptFile',
        content: content,
      });
      console.log('workspace script saved');
    });
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
