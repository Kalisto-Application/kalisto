import React, { useContext } from 'react';
import iconCopy from '../../assets/icons/copy.svg';
import { Context } from '../state';
import { CodeViewer } from '../ui/Editor';

export const ScriptResponse: React.FC = () => {
  const ctx = useContext(Context);

  return (
    <div className="w-1/2 bg-textBlockFill">
      <div className="p-[5.5px] text-right">
        <button>
          <img src={iconCopy} />
        </button>
      </div>
      <CodeViewer
        fileId={ctx.state.scriptResponse}
        text={ctx.state.scriptResponse}
      />
    </div>
  );
};
