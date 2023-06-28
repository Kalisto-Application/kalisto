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

  return(
    <div>
    <ControlledTreeEnvironment 
      items={itemsData}
      getItemTitle={item => item.data}
      viewState={{
        "tree1": {
          expandedItems: items.map(it => it.fullName),
          focusedItem: defaultFocusedItem,
        }
      }}
    >
        <Tree treeId="tree1" rootItem="root" treeLabel="methods collection" />
    </ControlledTreeEnvironment>
    </div>
    );
};
