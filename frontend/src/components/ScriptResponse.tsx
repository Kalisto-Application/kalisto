import React, { useContext } from 'react';
import { Context } from '../state';
import { CodeViewer } from '../ui/Editor';

export const ScriptResponse: React.FC = () => {
  const ctx = useContext(Context);

  return (
    <div className="w-1/2 bg-textBlockFill">
      <CodeViewer
        fileId={ctx.state.scriptResponse}
        text={ctx.state.scriptResponse}
      />
    </div>
  );
};
