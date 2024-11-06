import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useBusinessInfo } from "../hooks/useBusinessInfo.ts";
import { useApplicationInfo } from "../hooks/useApplicationInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";

export const ApplicationInfo = () => {
    const business_ids = useBusinessIds();
    const data = useApplicationInfo({ business_ids })[0];

    if (data != null) {
        return (
            <div className="content">
                <Card interactive={true} >
                    <H3>Account Information</H3>
                    <FormGroup label="Name"
                        labelFor="name" >
                        <InputGroup id="name" placeholder={data.id} />
                    </FormGroup>
                    <FormGroup label="Email"
                        labelFor="email" >
                        <InputGroup id="email" placeholder={"blah"} />
                    </FormGroup>
                    <Button>Save changes</Button>
                </Card>
            </div>
        );
    }
}