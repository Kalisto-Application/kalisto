import React, { useContext } from "react";
import { Context } from "../state";

import Text from "../ui/Text"; 

export const ApiError: React.FC = () => {
    const ctx = useContext(Context);

    return <Text text={ctx.state.apiError} props={{className: 'text-red'}}></Text>
}

export default ApiError
