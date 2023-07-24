import React, { useState, useEffect, useContext } from "react";
import { ControlledTreeEnvironment, Tree, TreeItemIndex, TreeItem } from 'react-complex-tree';
import { models } from "../../wailsjs/go/models";
import "react-complex-tree/lib/style-modern.css";
import { Context } from "../state";

interface MethodCollectionProps {
  services: models.Service[];
  selectedItem?: string;
};

const findMethod = (s: models.Service[], name: string): models.Method | undefined => {
  for (const service of s) {
    for (const method of service.methods) {
      if (method.fullName == name) {
        return method
      }
    }
  }
}

type Data = {
  display: string;
  isMethod: boolean;
}

export const MethodCollection: React.FC<MethodCollectionProps> = ({ services, selectedItem }) => {
  const ctx = useContext(Context);
  
  const serviceNames = services.map(it => it.fullName)
  const [expandedItems, setExpandedItems] = useState<TreeItemIndex[]>(serviceNames);

  useEffect(() => {
    setExpandedItems(serviceNames)
  }, [services])

  const itemsData: Record<TreeItemIndex, TreeItem<Data>> = {root: {
    index: 'root',
    isFolder: true,
    data: {display: 'Root item', isMethod: false},
    children: serviceNames,
  }}
  services.forEach(it => {
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
          const method = findMethod(services, item.index as string)

          ctx.dispatch({type: 'activeMethod', activeMethod: method!})
        }
        setExpandedItems([...expandedItems, item.index])}
      }
      onCollapseItem={item => {
        setExpandedItems(expandedItems.filter(expandedItemIndex => expandedItemIndex !== item.index))
      }}
      onFocusItem={item => {
        if (item.data.isMethod) {
          const method = findMethod(services, item.index as string)
          ctx.dispatch({type: 'activeMethod', activeMethod: method!})
        }
      }}
    >
        <Tree treeId="1" rootItem="root" treeLabel="methods collection" />
    </ControlledTreeEnvironment>
    </div>
    );
};
