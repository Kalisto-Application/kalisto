import React, { useContext } from "react";
import { Context } from "../state";

import Text from "../ui/Text"; 

export const ScriptError: React.FC = () => {
    const ctx = useContext(Context);

    return <Text text={ctx.state.scriptError} props={{className: 'text-red'}}></Text>
}

export default ScriptError
