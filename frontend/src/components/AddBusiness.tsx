import { Button, Card, FormGroup, H3, H5, Icon, InputGroup, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useBusinessInfo } from "../hooks/useBusinessInfo.ts";
import { NewBusinessInfo, useCreateBusiness } from "../hooks/useCreateBusiness.ts";


interface AddBusinessProps {
    setBusinessAdd: (boolean) => void;
}

export const AddBusiness = ({ setBusinessAdd }: AddBusinessProps) => {

    const { createBusiness } = useCreateBusiness();

    const [userInfo, setUserInfo] = useState<NewBusinessInfo>({ name: "", desc: "", website: "" });

    const handleChangeInfo = (id, value) => {
        setUserInfo((prevInfo) => ({
            ...prevInfo,
            [id]: value
        }))
    }

    const handleSubmit = () => {
        createBusiness(userInfo);
        setBusinessAdd(false);
    }

    const handleCancel = () => {
        setBusinessAdd(false);
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', padding: '20px' }}>
            <Card interactive={true} >
                <FormGroup
                    label="Business name"
                    labelFor="name"
                >
                    <InputGroup id="name" value={userInfo.name} onValueChange={(value) => handleChangeInfo("name", value)} placeholder="Business Name" />
                </FormGroup>
                <FormGroup
                    label="Description"
                    labelFor="desc"
                >
                    <InputGroup id="desc" value={userInfo.desc} onValueChange={(value) => handleChangeInfo("desc", value)} placeholder="Description" />
                </FormGroup>
                <FormGroup
                    label="Website"
                    labelFor="website"
                >
                    <InputGroup id="website" value={userInfo.website} onValueChange={(value) => handleChangeInfo("website", value)} placeholder="Website" />
                </FormGroup>
                <Button onClick={() => handleSubmit()}>Submit</Button>
                <Button onClick={() => handleCancel()}>Cancel</Button>
            </Card>
        </div>
    );
}
