import {
    Button,
    Card,
    FormGroup,
    H3,
    InputGroup,
    Intent,
    Classes,
    Callout,
    HTMLSelect,
    Spinner
} from "@blueprintjs/core";
import React, { useState } from "react";
import { NewPostingInfo, useCreatePosting } from "../hooks/useCreatePosting.ts";

interface AddPostingProps {
    setPostingAdd: (value: boolean) => void;
    fetchData: () => Promise<void>;
    businesses: string[];
}

interface FormErrors {
    title?: string;
    desc?: string;
    pay?: string;
    business_id?: string;
}

export const AddPosting = ({ setPostingAdd, fetchData, businesses }: AddPostingProps) => {
    const { createPosting } = useCreatePosting();
    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState<FormErrors>({});

    const [postingInfo, setPostingInfo] = useState<NewPostingInfo>({
        title: "",
        desc: "",
        pay: 0,
        time_est: 1,
        business_id: businesses[0] || ""
    });

    const validateForm = (): boolean => {
        const newErrors: FormErrors = {};

        if (!postingInfo.title.trim()) {
            newErrors.title = "Title is required";
        }

        if (!postingInfo.desc.trim()) {
            newErrors.desc = "Description is required";
        }

        if (postingInfo.pay < 0) {
            newErrors.pay = "Pay cannot be negative";
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleChangeInfo = (id: keyof NewPostingInfo, value: string) => {
        setPostingInfo((prevInfo) => ({
            ...prevInfo,
            [id]: value
        }));

        if (errors[id]) {
            setErrors((prev) => ({
                ...prev,
                [id]: undefined
            }));
        }
    };

    const handleChangeNumber = (id: keyof NewPostingInfo, value: string) => {
        let tempValue = 0;
        if (!isNaN(Number(value))) {
            tempValue = Number(value);
        }

        setPostingInfo((prevInfo) => ({
            ...prevInfo,
            [id]: tempValue
        }));
    };

    const handleSubmit = async () => {
        if (!validateForm()) {
            return;
        }

        try {
            setIsLoading(true);
            await createPosting(postingInfo);
            await fetchData();
            setPostingAdd(false);
        } catch (error) {
            console.error("Failed to create posting:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="content">
            <Card
                style={{
                    width: "100%",
                    maxWidth: "500px",
                    margin: "auto",
                    padding: "2rem"
                }}
            >
                <H3>Add New Posting</H3>

                <div style={{ marginTop: "1.5rem" }}>
                    {Object.values(errors).some(error => error) && (
                        <Callout
                            intent={Intent.DANGER}
                            icon="error"
                            style={{ marginBottom: "1rem" }}
                        >
                            Please fix the errors below to continue
                        </Callout>
                    )}

                    <FormGroup
                        label="Title"
                        labelFor="title"
                        intent={errors.title ? Intent.DANGER : Intent.NONE}
                        helperText={errors.title}
                    >
                        <InputGroup
                            id="title"
                            value={postingInfo.title}
                            onChange={(e) => handleChangeInfo("title", e.target.value)}
                            placeholder="Enter posting title"
                            intent={errors.title ? Intent.DANGER : Intent.NONE}
                        />
                    </FormGroup>

                    <FormGroup
                        label="Description"
                        labelFor="desc"
                        intent={errors.desc ? Intent.DANGER : Intent.NONE}
                        helperText={errors.desc}
                    >
                        <InputGroup
                            id="desc"
                            value={postingInfo.desc}
                            onChange={(e) => handleChangeInfo("desc", e.target.value)}
                            placeholder="Enter description"
                            intent={errors.desc ? Intent.DANGER : Intent.NONE}
                        />
                    </FormGroup>

                    <FormGroup
                        label="Pay"
                        labelFor="pay"
                        intent={errors.pay ? Intent.DANGER : Intent.NONE}
                        helperText={errors.pay}
                    >
                        <InputGroup
                            id="pay"
                            value={postingInfo.pay.toString()}
                            onChange={(e) => handleChangeNumber("pay", e.target.value)}
                            placeholder="Enter pay amount"
                            intent={errors.pay ? Intent.DANGER : Intent.NONE}
                            type="number"
                        />
                    </FormGroup>

                    <FormGroup
                        label="Business"
                        labelFor="business_id"
                    >
                        <HTMLSelect
                            id="business_id"
                            value={postingInfo.business_id}
                            onChange={(e) => handleChangeInfo("business_id", e.target.value)}
                            options={businesses}
                            fill={true}
                        />
                    </FormGroup>

                    <div style={{
                        display: "flex",
                        justifyContent: "flex-end",
                        gap: "0.75rem",
                        marginTop: "2rem"
                    }}>
                        <Button
                            onClick={() => setPostingAdd(false)}
                            minimal={true}
                        >
                            Cancel
                        </Button>
                        <Button
                            onClick={handleSubmit}
                            intent={Intent.PRIMARY}
                            disabled={isLoading}
                        >
                            {isLoading ? <Spinner size={16} /> : "Submit"}
                        </Button>
                    </div>
                </div>
            </Card>
        </div>
    );
};