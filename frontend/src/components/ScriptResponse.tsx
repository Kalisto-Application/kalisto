import React, { useContext } from 'react';
import { CodeViewer } from '../ui/Editor';
import { Context } from '../state';

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
