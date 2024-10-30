import { Card } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo";

export const Account = () => {
    const { data } = useAccountInfo();
    console.log(data);
    return (
        <div>


        </div>
    );
}

export default Account;