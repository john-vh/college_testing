import { Button, Card, FormGroup, H3, H5, Icon, InputGroup, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { usePostingInfo } from "../hooks/usePostingInfo.ts";
import { NewPostingInfo, useCreatePosting } from "../hooks/useCreatePosting.ts";
import { Select } from "@blueprintjs/select";


interface AddPostingProps {
    setPostingAdd: (boolean) => void;
    fetchData: () => Promise<void>;
    businesses: string[]
}

export const AddPosting = ({ setPostingAdd, fetchData, businesses }: AddPostingProps) => {
    const { createPosting } = useCreatePosting();
    const [postingInfo, setPostingInfo] = useState<NewPostingInfo>({ title: "", desc: "", pay: 0, time_est: 1, business_id: businesses[0] });

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const handleChangeInfo = (id, value) => {
        setPostingInfo((prevInfo) => ({
            ...prevInfo,
            [id]: value
        }))
    }

    const handleChangeNumber = (id, value) => {
        let tempValue = 0;
        if (!isNaN(Number(value))) {
            tempValue = Number(value);
        }

        setPostingInfo((prevInfo) => ({
            ...prevInfo,
            [id]: tempValue
        }))
    }

    const handleSubmit = async () => {
        createPosting(postingInfo);
        await delay(2000);
        await fetchData();
        setPostingAdd(false);
    }

    const handleCancel = () => {
        setPostingAdd(false);
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', padding: '20px', marginLeft: "200px" }}>
            <Card interactive={true} >
                <FormGroup
                    label="Posting title"
                    labelFor="title"
                >
                    <InputGroup id="title" defaultValue={postingInfo.title} onValueChange={(value) => handleChangeInfo("title", value)} placeholder="Title" />
                </FormGroup>
                <FormGroup
                    label="Description"
                    labelFor="desc"
                >
                    <InputGroup id="desc" defaultValue={postingInfo.desc} onValueChange={(value) => handleChangeInfo("desc", value)} placeholder="Description" />
                </FormGroup>
                <FormGroup
                    label="Pay"
                    labelFor="pay"
                >
                    <InputGroup id="pay" defaultValue={postingInfo.pay.toString()} onValueChange={(value) => handleChangeNumber("pay", value)} placeholder="Pay" />
                </FormGroup>
                <FormGroup
                    label="Business"
                    labelFor="business"
                >
                    <InputGroup id="business" defaultValue={postingInfo.business_id} onValueChange={(value) => handleChangeInfo("business_id", value)} placeholder="Business" />
                </FormGroup>

                <Button onClick={() => handleSubmit()}>Submit</Button>
                <Button onClick={() => handleCancel()}>Cancel</Button>
            </Card>
        </div >
    );
}
