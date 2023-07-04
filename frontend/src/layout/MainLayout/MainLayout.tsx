import React, {useContext, useEffect, useState} from "react";

import Header from "../Header";
import Sidebar from "../Sidebar";

import { WorkspaceList, WorkspaceItem } from "../../components/Workspaces";
import { FindWorkspaces, GetWorkspace, NewWorkspace } from "../../../wailsjs/go/api/Api";

import { MethodItem, ServiceItem } from "../../components/MethodCollectionView";
import { MainPageContent } from "../../pages/MainPageContent";
import { Context } from "../../state";

export const MainLayout: React.FC = () => {;
  const [workspaceId, setWorkspaceId] = useState<string>('');
  const [workspaceItems, setWorkspaceItems] = useState<WorkspaceItem[]>([]);
  const [methodItems, setMethodItems] = useState<ServiceItem[]>([]);
  const [method, setMethod] = useState<MethodItem>()
  const ctx = useContext(Context);

  useEffect(() => {
    FindWorkspaces()
    .then(res => {
      if (res == null) {
        return
      }
      const items = res.map(it => ({id: it.id, name: it.name}))
      if (items.length == 0) {
        return
      }

      let latest = res[0]
      res.forEach(it => {
        if (it.lastUsage > latest.lastUsage) {
          latest = it.lastUsage
        }
      })

      const getFirstMethod = (): MethodItem | undefined => {
        for (const ws of res) {
          for (const service of ws.spec.services) {
            for (const m of service.methods) {
              return {name: m.name, fullName: m.fullName, requestExample: m.requestExample}
            }
          }
        }
      }
      const fristMethod = getFirstMethod()

      setWorkspaceId(latest.id)
      setWorkspaceItems(items.map(it => ({id: it.id, name: it.name, active: it.id == latest.id})))
      setMethodItems(res.find(it => it.id == latest.id)!.spec.services.map(s => ({name: s.name, fullName: s.fullName, methods: s.methods.map(met => ({name: met.name, fullName: met.fullName, requestExample: met.requestExample}))})));
      setMethod(fristMethod);
      ctx.dispatch({type: 'changeRequestText', text: fristMethod?.requestExample || ''})
    })
    .catch(err => console.log('error on find workspaces: ', err))
  }, [])

  const setActiveWorkspace = (id: string) => {
    setWorkspaceId(id);
    setWorkspaceItems(items => items.map(it => ({id: it.id, name: it.name, active: it.id == id})))
    GetWorkspace(id).then(res => {
      setMethodItems(res.spec.services.map(s => ({name: s.name, fullName: s.fullName, methods: s.methods.map(met => ({name: met.name, fullName: met.fullName, requestExample: met.requestExample}))})));
    }).catch(err => console.log(`error on get workspace by id==${id}: `, err))
  }

  const newWorkspace = () => {
    NewWorkspace()
    .then(res => {
      setWorkspaceId(res.id)
      setWorkspaceItems(prev => {
        const items = prev.map(it => ({id: it.id, name: it.name, active: false}))
        items.push({id: res.id, name: res.name, active: true})
        return items;
      });
      setMethodItems(res.spec.services.map(s => ({name: s.name, fullName: s.fullName, methods: s.methods.map(met => ({name: met.name, fullName: met.fullName, requestExample: met.requestExample}))})));
    })
    .catch(err => console.log('error on new workspace: ', err))
  }

  return (<div className="flex h-screen flex-col">
    <Header>
      <WorkspaceList items={workspaceItems} setActive={setActiveWorkspace} newWorkspace={newWorkspace}/>
    </Header>
    <div className="flex flex-1">
      <Sidebar>
      </Sidebar>
      <main className="flex-1">
        <MainPageContent workspaceId={workspaceId} methodItems={methodItems} method={method} setActiveMethod={setMethod} inputText={ctx.state.requestText}/>
      </main>
    </div>
  </div>);
}
