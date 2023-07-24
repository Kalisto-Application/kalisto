import React, {useContext, useEffect} from "react";

import Header from "../Header";
import Sidebar from "../Sidebar";

import { WorkspaceList } from "../../components/Workspaces";
import { ApiPage } from "../../pages/ApiPage";
import { Context } from "../../state";
import { FindWorkspaces } from "../../../wailsjs/go/api/Api";
import { models } from "../../../wailsjs/go/models";

export const MainLayout: React.FC = () => {;
  const ctx = useContext(Context);

  useEffect(() => {
    FindWorkspaces()
    .then(res => {
      if (res == null) {
        return
      }

      let latest = res[0]
      res.forEach(it => {
        if (it.lastUsage > latest.lastUsage) {
          latest = it.lastUsage
        }
      })

      const getFirstMethod = (): models.Method | undefined => {
        for (const ws of res) {
          for (const service of ws.spec.services) {
            for (const m of service.methods) {
              return m
            }
          }
        }
      }
      const fristMethod = getFirstMethod()
      if (fristMethod) {
        ctx.dispatch({type: 'activeMethod', activeMethod: fristMethod});
      }
      ctx.dispatch({type: 'workspaceList', workspaceList: res});
    })
    .catch(err => console.log('error on find workspaces: ', err))
  }, [])

  return (<div className="flex h-screen flex-col">
    <Header>
      <WorkspaceList items={ctx.state.workspaceList} activeWorkspace={ctx.state.activeWorkspace} />
    </Header>
    <div className="flex flex-1">
      <Sidebar />
      <main className="flex-1">
        <ApiPage workspace={ctx.state.activeWorkspace} method={ctx.state.activeMethod} inputText={ctx.state.requestText} />
      </main>
    </div>
  </div>);
}
