import { Button, Card, FormGroup, H3, H5, Icon, InputGroup, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { usePostingInfo } from "../hooks/usePostingInfo.ts";
import { NewPostingInfo, useCreatePosting } from "../hooks/useCreatePosting.ts";
import { Select } from "@blueprintjs/select";


interface AddPostingProps {
    setPostingAdd: (boolean) => void;
    businesses: Set<string>
}

export const AddPosting = ({ setPostingAdd, businesses }: AddPostingProps) => {
    const { createPosting } = useCreatePosting();
    const [postingInfo, setPostingInfo] = useState<NewPostingInfo>({ title: "", desc: "", pay: 0, time_est: 1, business_id: businesses[0] });

    const handleChangeInfo = (id, value) => {
        setPostingInfo((prevInfo) => ({
            ...prevInfo,
            [id]: value
        }))
    }

    const handleSubmit = () => {
        createPosting(postingInfo);
        setPostingAdd(false);
    }

    const handleCancel = () => {
        setPostingAdd(false);
    }

    return (
        <div style={{ display: 'flex', flexDirection: 'column', padding: '20px' }}>
            <Card interactive={true} >
                <FormGroup
                    label="Posting title"
                    labelFor="title"
                >
                    <InputGroup id="title" value={postingInfo.title} onValueChange={(value) => handleChangeInfo("title", value)} placeholder="Posting Title" />
                </FormGroup>
                <FormGroup
                    label="Description"
                    labelFor="desc"
                >
                    <InputGroup id="desc" value={postingInfo.desc} onValueChange={(value) => handleChangeInfo("desc", value)} placeholder="Description" />
                </FormGroup>
                <FormGroup
                    label="Business"
                    labelFor="business"
                >
                    <InputGroup id="business" value={postingInfo.business_id} onValueChange={(value) => handleChangeInfo("business_id", value)} placeholder="Business" />
                </FormGroup>

                <Button onClick={() => handleSubmit()}>Submit</Button>
                <Button onClick={() => handleCancel()}>Cancel</Button>
            </Card>
        </div >
    );
}
