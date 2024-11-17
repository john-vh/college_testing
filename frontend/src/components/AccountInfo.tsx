import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";

export const AccountInfo = () => {

    const data = useAccountInfo();

    if (data != null) {
        return (
            <div className="content">
                <Card interactive={true} >
                    <H3>Account Information</H3>
                    <FormGroup label="Name"
                        labelFor="name" >
                        <InputGroup id="name" defaultValue={data.name} />
                    </FormGroup>
                    <FormGroup label="Email"
                        labelFor="email" >
                        <InputGroup id="email" defaultValue={data.email} />
                    </FormGroup>
                    <Button>Save changes</Button>
                </Card>
            </div>
        );
    }
}

export default AccountInfo;