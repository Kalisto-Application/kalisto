import React, {SyntheticEvent, useState, useEffect} from "react";
import { ControlledTreeEnvironment, UncontrolledTreeEnvironment, Tree, StaticTreeDataProvider, TreeDataProvider, TreeItemIndex, TreeItem, Disposable, ExplicitDataSource } from 'react-complex-tree';
import { EventEmitter } from "react-complex-tree/src/EventEmitter";
import "react-complex-tree/lib/style-modern.css";

interface MethodCollectionProps {
  onClick: (e: SyntheticEvent, data: string) => void;
  items: MethodItem[];
};

export type MethodItem = {
  name: string;
  fullName: string;
  methods : {
    name: string;
    fullName: string;
  }[];
}

type TreeItemRecord = Record<TreeItemIndex, TreeItem<string>>

class DataProvider<TreeItemRecord> implements TreeDataProvider {
  private data: ExplicitDataSource;
  private onDidChangeTreeDataEmitter = new EventEmitter<TreeItemIndex[]>();

  constructor(
    items: Record<TreeItemIndex, TreeItem<TreeItemRecord>>
  ) {
    this.data = { items };
  }

  public async getTreeItem(itemId: TreeItemIndex): Promise<TreeItem> {
    return this.data.items[itemId];
  }

  public async onChangeItemChildren(
    itemId: TreeItemIndex,
    newChildren: TreeItemIndex[]
  ): Promise<void> {
    this.data.items[itemId].children = newChildren;
    this.onDidChangeTreeDataEmitter.emit([itemId]);
  }

  public onDidChangeTreeData(
    listener: (changedItemIds: TreeItemIndex[]) => void
  ): Disposable {
    const handlerId = this.onDidChangeTreeDataEmitter.on(payload =>
      listener(payload)
    );
    return { dispose: () => this.onDidChangeTreeDataEmitter.off(handlerId) };
  }

  public async onRenameItem(item: TreeItem<TreeItemRecord>, name: string): Promise<void> {
      this.data.items[item.index].data = name;
      this.onDidChangeTreeDataEmitter.emit([item.index]);
  }
}

export const MethodCollection: React.FC<MethodCollectionProps> = ({ onClick, items }) => {
  const itemsData: Record<TreeItemIndex, TreeItem<string>> = {root: {
    index: 'root',
    isFolder: true,
    data: 'Root item',
    children: items.map(it => it.fullName),
  }}
  items.forEach(it => {
    itemsData[it.fullName] = {
      index: it.fullName,
      data: it.name,
      isFolder: true,
      children: it.methods.map(met => met.fullName),
    }
    it.methods.forEach(met => {
      itemsData[met.fullName] = {
        index: met.fullName,
        data: met.name,
      }
    })
  })

  let defaultFocusedItem: TreeItemIndex | undefined;
  if (items.length > 0 && items[0].methods.length > 0) {
    defaultFocusedItem = items[0].methods[0].fullName
  }

  console.log('tree: ', itemsData);

  return(
    <div>
      <UncontrolledTreeEnvironment dataProvider={new DataProvider(itemsData)} getItemTitle={item => item.data} viewState={{}}>
      {/* <UncontrolledTreeEnvironment dataProvider={new StaticTreeDataProvider(itemsData, (item, data) => ({ ...item, data }))} getItemTitle={item => item.data} viewState={{}}> */}
    {/* <ControlledTreeEnvironment 
      items={itemsData}
      getItemTitle={item => item.data}
      viewState={{
        "tree1": {
          expandedItems: items.map(it => it.fullName),
          focusedItem: defaultFocusedItem,
        }
      }}
    > */}
        <Tree treeId="tree1" rootItem="root" treeLabel="methods collection" />
      </UncontrolledTreeEnvironment>
    {/* </ControlledTreeEnvironment> */}
    </div>
    );
};
