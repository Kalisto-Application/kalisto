import React from 'react';
import SearchBox from '../ui/SearchBox';
import { WorkspaceList } from '../components/Workspaces';
import { ScriptEditor } from '../components/ScriptEditor';
import { ScriptResponse } from '../components/ScriptResponse';
import TabList from '../ui/TabList';
import ScriptError from '../components/ScriptError';
import { ScriptSender } from '../components/ScriptSender';
import ScriptCollectionView from '../components/ScriptCollectionView';

export const ScriptingPage: React.FC = () => {
  return (
    <div className="flex w-full flex-1">
      <div className="flex flex-[0_0_220px] flex-col justify-items-start">
        <WorkspaceList />
        <SearchBox />
        <ScriptCollectionView />
      </div>
      <div className="flex flex-1 flex-col">
        <TabList />
        <ScriptSender />
        <div className="flex flex-1">
          <ScriptEditor />
          <ScriptResponse />
        </div>
        <ScriptError />
      </div>
    </div>
  );
};
