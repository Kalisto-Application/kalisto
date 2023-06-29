import React, {SyntheticEvent, useState, useEffect} from "react";
import { ControlledTreeEnvironment, UncontrolledTreeEnvironment, Tree, StaticTreeDataProvider, TreeDataProvider, TreeItemIndex, TreeItem, Disposable, ExplicitDataSource } from 'react-complex-tree';
import { EventEmitter } from "react-complex-tree/src/EventEmitter";
import "react-complex-tree/lib/style-modern.css";

interface MethodCollectionProps {
  setActiveMethod: (fullName: string) => void;
  items: MethodItem[];
  defaultFocused: string;
};

export type MethodItem = {
  name: string;
  fullName: string;
  methods : {
    name: string;
    fullName: string;
  }[];
}

type Data = {
  display: string;
  isMethod: boolean;
}

export const MethodCollection: React.FC<MethodCollectionProps> = ({ setActiveMethod, items, defaultFocused }) => {
  const serviceNames = items.map(it => it.fullName)
  const [expandedItems, setExpandedItems] = useState<TreeItemIndex[]>(serviceNames);
  const [selectedItem, setSelectedItem] = useState<TreeItemIndex | undefined>(defaultFocused);

  useEffect(() => {
    setExpandedItems(serviceNames)
  }, [items])

  const itemsData: Record<TreeItemIndex, TreeItem<Data>> = {root: {
    index: 'root',
    isFolder: true,
    data: {display: 'Root item', isMethod: false},
    children: serviceNames,
  }}
  items.forEach(it => {
    itemsData[it.fullName] = {
      index: it.fullName,
      data: {display: it.name, isMethod: false},
      isFolder: true,
      children: it.methods.map(met => met.fullName),
    }
    it.methods.forEach(met => {
      itemsData[met.fullName] = {
        index: met.fullName,
        data: {display: met.name, isMethod: true},
      }
    })
  })

  return(
    <div>
    <ControlledTreeEnvironment 
      items={itemsData}
      getItemTitle={item => item.data.display}
      viewState={{
        "1": {
          expandedItems: expandedItems,
          focusedItem: selectedItem,
        }
      }}
      onExpandItem={item => {
        if (item.data.isMethod) {
          setActiveMethod(item.index as string);
        }
        setExpandedItems([...expandedItems, item.index])}
      }
      onCollapseItem={item => {
        setExpandedItems(expandedItems.filter(expandedItemIndex => expandedItemIndex !== item.index))
      }}
      onFocusItem={item => {
        if (item.data.isMethod) {
          setActiveMethod(item.index as string);
        }
        setSelectedItem(item.index)
      }}
    >
        <Tree treeId="1" rootItem="root" treeLabel="methods collection" />
    </ControlledTreeEnvironment>
    </div>
    );
};
